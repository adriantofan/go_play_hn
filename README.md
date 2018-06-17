databse loading when compiled optimized  3.938896335 
databse loading when compiled optimized after removing year aggregation 3.443510763 

## Notes for improvement:

* The deterministic algorithm cannot be easily scaled. The solution seem to be non deterministic algorithms such as hyperLogLog and family. Basicaly data processing can be distributed to worker nodes and results merged centraly as follows :

    - grouping by url - long tail problem, a few machines will process the most popular urls. I cannot see how to evenly distribute popular urls.

    - grouping by date , let's say days. The problem is that the workers need to return to the server the complete hash map for the day and then do a central merge. If not done the agregate computation of top n can be (very) wrong. Imagine top 1 of node 1: {m:5,x:4} node 2 {b:6,x:4}. Agregate of top 1 is {b:6,m:5 }. In reality the result is {x:8, b:6}


* Given that AddLog is called on successive logs it might be possible to optimize it further such that we don't lookup the same log chain at each step

* It doesn't seem realistic to work with more than a few days. Given that logs are contignous the hash of hash trie implementation could easily be replaced with an big array reducing further startup and look up time.

* it might be posible to not work with urls and replace them with a data structure that is basicaly a string hash and byte range in the file. Given that the file is memmory mapped, it can work nicely, but eventualy it will be verry disk bound.

* date parsing and conversion to components can be much improved given the standard format 

