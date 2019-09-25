## Understand the differences between shared memory data-parallelism and distributed data-parallelism

### Shared memory data-parallelism 

Let's say you have like a very...normal... computer, like this one (this is a computer, trust me) :

<p align="center">
  <img width="460" height="300" src="https://4.imimg.com/data4/RQ/PS/MY-25091456/how-to-donate-computer-1-500x500.jpg">
</p>

And thanks to this super amazing computer, you'd like to do some ETL-stylish computations on a logs file that contains all the events recorded by your online mobile game that day.
This logs file is going to be a **collection** of data. On this collection, you'd like to apply some transformations such as filtering records, making key/value pairs and so on...

So in case of a shared memory data-parallelism, in order to do these computations, the program has to split up this collection of data into several chunks via some kind of partionning algorithm. You must understand that this collection of data is sitting on our computer. Nowhere else.

After that, it has to spin up a few tasks or some kind of workers/threads that will independently operate on the chunks they received, **in parallel**. And after that, they will combine their result into one output.

[... Work in progress ...]

```golang
// Making a PCollection<record> of that logs file
logs := textio.Read(s, "../path/to/logs/file")
```



