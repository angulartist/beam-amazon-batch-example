## What is shuffling in the distributed data processing world?

### The dataset

Let's say we've got a super simple logs file containing millions of financial transactions, and we would like to **count the occurrences of each type of credit card** in order to build some analytics (how our customer consume etc.).

| id  | sender  |  receiver |  card_type | amount |
|---|---|---|---|---|
| 1  | Bob  |  Alan |  Visa |  42.00 |
| 2  |  Alan | Dylan  | Mastercard  |  69.40 |
| 3  | Lucie  | Eric  |  Visa |  11.90 |

### The map phase

First of, we have to build our PCollection of transactions from this source text file :

```golang
// Representation of a transaction
type Transaction struct {
    Id int
    Sender string
    Receiver string
    CardType string
    Amount float64
}

// NOTHING -> PCollection<Transaction>
transactions := textio.Read(s, '../path/to/file.txt')
```

Then, we must extract the credit card type from each transaction and build a key-value pair where the key is going to be the card type and the value's gonna be 1. This approach will allow us to run a reduce function in order to aggregate the sum.

```golang
// NOTHING -> PCollection<Transaction>
transactions := textio.Read(s, '../path/to/file.txt')

// New
// PCollection<Transaction> -> PCollection<KV<string, int>>
pairedWithOne := beam.ParDo(s, func (tsx Transaction) (string, int) {
    return tsx.CardType, 1
}, transactions)
```

**What might this look like on the cluster?**

Let's say we have a cluster of 3 worker nodes. Each node has a chunk (or partition) of our data and the user code (UDFs) to apply to each element. We will get a bunch of key-value pairs on each node, just like that :

| Worker1  | Worker2  | Worker3  |
|---|---|---|
| (Visa, 1)  | (Mastercard, 1)  | (Visa, 1)  |
| (Mastercard, 1)  | (Visa, 1)  | (Mastercard, 1)  |
| (Mastercard, 1)  | (Visa, 1)  | (Visa, 1)  |
| (Visa, 1)  | (Maestro, 1)  | (Mastercard, 1)  |
| (Visa, 1)  | (Mastercard, 1)  |  (Maestro, 1)  |
| (Maestro, 1)  | (Mastercard, 1)  | (Visa, 1)  |

### The reduce phase

Now, we want to group all values with the same key (the credit card Type) together in order to aggregate a sum. So we'll use a combination of a `GroupByKey` with a reducer.

```golang
// NOTHING -> PCollection<Transaction>
transactions := textio.Read(s, '../path/to/file.txt')

// PCollection<Transaction> -> PCollection<KV<string, int>>
pairedWithOne := beam.ParDo(s, func (tsx Transaction) (string, int) {
    return tsx.CardType, 1
}, transactions)

// New
// PCollection<KV<string, int>> -> PCollection<GBK<string, iter<int>>>
grouppedByKey := beam.GroupByKey(s, pairedWithOne)

// PCollection<GBK<string, iter<int>>> -> PCollection<KV<string, int>>
counted := beam.ParDo(s, func(key string, iter []int) (key, int) {
    return key, sum(iter)
}, grouppedByKey)
```

> Note: GroupByKey is a PTransform that takes a PCollection of type KV<A,B>, groups the values by key and windows, and returns a PCollection of type GBK<A,B> representing a map from each distinct key and window of the input PCollection to an iterable over all the values associated with that key in the input per window. Each key in the output PCollection is unique within each window.

**What might this look like on the cluster?**

Right before the reduce phase happens, key-value pairs will **move across the network** so the same keys will be gathered on the same machine. In this example, we assume that each worker node is hosting a unique key. Each key has an iterable of integers as a value, which represents the number of time we've seen a credit card type.

*GroupByKey*

| Worker1  | Worker2  | Worker3  |
|---|---|---|
| (Visa, [1, 1, 1, 1, 1, 1, 1, 1])  | (Mastercard, [1, 1, 1, 1, 1, 1, 1])  | (Maestro, [1, 1, 1])  |

*Reduction*

| Worker1  | Worker2  | Worker3  |
|---|---|---|
| (Visa, 8)  | (Mastercard, 7)  | (Maestro, 3)  |

### Houston, we have a problem

But didn't we talk about shuffling? Just before the reduce phase, when data is moving across the network, this is called **shuffling**. And when shuffling a small dataset might be *okay*... on a huge dataset this will introduce **latency** because too much network communication kills performance. 

So you **do not want sending all of your data across the network if it's not absolutely required**.

(don't do it, for real)

So, in our case, the `GroupByKey` was a naive approach. We have elements of the same type and we want to calculate a sum, hum, can we do it in a more efficient way? Can we reduce before doing the shuffle to reduce the data that will be sent over the network? Yay!

![](https://media.makeameme.org/created/yay.jpg)

### SumPerKey/CombinePerKey

`SumPerKey` or `CombinePerKey` are a combination of first doing a GroupByKey and then **reduce-ing** on all the values grouped by that key. This is way more efficient than using each separately. And we can use them because we are aggregating values of the same type. So that's perfectly fine.

*Using GroupByKey*

```golang
// PCollection<KV<string, int>> -> PCollection<GBK<string, iter<int>>>
grouppedByKey := beam.GroupByKey(s, pairedWithOne)

// PCollection<GBK<string, iter<int>>> -> PCollection<KV<string, int>>
counted := beam.ParDo(s, func(key string, iter []int) (key, int) {
    return key, sum(iter)
}, grouppedByKey)
```

*Using SumPerKey*

```golang
// PCollection<KV<string, int>> -> PCollection<KV<string, int>>
stats.SumPerKey(s, pairedWithOne)
```

**What might this look like on the cluster?**

Because `SumPerKey` is an associative reduction : magic happens! In fact, it will reduces the data on the mapper side first before sending the aggregated results to the down stream. And because results are already pre-aggregated, the data that will be sent over the network for the final reduction will be greatly limited. So we get the same output but **faster**.

![Magic](https://media.makeameme.org/created/magic-its.jpg)

*Mappers will build key-value pairs*

| Worker1  | Worker2  | Worker3  |
|---|---|---|
| (Visa, 1)  | (Mastercard, 1)  | (Visa, 1)  |
| (Mastercard, 1)  | (Visa, 1)  | (Mastercard, 1)  |
| (Mastercard, 1)  | (Visa, 1)  | (Visa, 1)  |
| (Visa, 1)  | (Maestro, 1)  | (Mastercard, 1)  |
| (Visa, 1)  | (Mastercard, 1)  |  (Maestro, 1)  |
| (Maestro, 1)  | (Mastercard, 1)  | (Visa, 1)  |

*Mappers will group elements by key*

| Worker1  | Worker2  | Worker3  |
|---|---|---|
| (Visa, [1, 1, 1])  | (Mastercard, [1, 1, 1])  | (Visa, [1, 1, 1])  |
| (Mastercard, [1, 1])  | (Visa, [1, 1])  | (Mastercard, [1, 1])  |
| (Maestro, [1])  | (Maestro, [1])  | (Maestro, [1]))  |

*Mappers will reduce the data they own locally*

| Worker1  | Worker2  | Worker3  |
|---|---|---|
| (Visa, 3)  | (Mastercard, 3)  | (Visa, 3)  |
| (Mastercard, 2)  | (Visa, 2)  | (Mastercard, 2)  |
| (Maestro, 1)  | (Maestro, 1)  | (Maestro, 1))  |

And then, we have the shuffle phase followed by the reduce phase.

*Aggregated output*

| Worker1  | Worker2  | Worker3  |
|---|---|---|
| (Visa, 8)  | (Mastercard, 7)  | (Maestro, 3)  |

### To sum up

Not all problems that can be solved by `GroupByKey` can be calculated with `SumPerKey` / `CombinePerKey`, in fact they require combining all your values into another value with the **exact same type**.

But just keep in mind that you can greatly optimise your data processing pipeline by reducing the shuffle as much as possible.

Happy coding. :poop:

