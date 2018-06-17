## What it does ?

Loads and processes data from a TSV file listing all queries performed on HN Search during a few days.

It then launches a web server that answer the following queries :

* `GET /1/queries/count/<DATE_PREFIX>`: returns a JSON object specifying the number of distinct queries that have been done during a specific time range
* `GET /1/queries/popular/<DATE_PREFIX>?size=<SIZE>`: returns a JSON object listing the top <SIZE> popular queries that have been done during a specific time range

## Architecture

The app loads logs one by one and adds them in to a trie-like structure. The trie organizes the logs by date.  For example given the logs :

        2015-08-01 00:03:42	term
        2015-08-01 00:03:43	term1    
        2015-08-01 00:03:44	term1    


is stored in a six level deep trie as

        2015                                <- aggregate for 2015
        |
        -> 08                               <- aggregate for 08 month of 2015
            |
            -> 01                           <- aggregate for 01 day of 08 month of 2015
                |
                -> 00                       <- aggregate for 00 hour for 2015-08-01
                    |
                    ->03                    <- aggregate for 03 minute of 2015-08-01 00:
                        |
                        -> 42               <- aggregate for 42 second of 2015-08-01 00:03

During loading, at each level, a aggregate term count is kept in a hash table. 

When all data is loaded, in a final processing step term counts are extracted form the hash table and kept sorted in a array as (term, count) pairs.

The trie is uses a hash table to access children nodes.

To answer how many distinct queries at 2015-08-01, just get the length of the  sorted (term, count) array at a third level 2015->08->01.

To answer what are the top 5 queries at 2015-08-01, just get first 5 elements of the sorted (term, count) array at a third level 2015->08->01.

Performance (amortized):
* Answer time is proportional to O(1) for both queries. 
* The startup time is around 3-4 seconds.
* The load time is N log N
* The memory usage is also proportional to N
* The constant factors in the previous estimates are quite high. Improvements are possible, keeping more or less the same algorithm. 


## Notes for improvement:

* The deterministic algorithm cannot be easily scaled. In order to scale up, data processing can be distributed to worker nodes and results merged centrally as follows :
    - grouping by url - long tail problem, a few machines will process the most popular urls. I cannot see how to evenly distribute popular urls.
    - grouping by date , let's say days. The problem is that the workers need to return to the server the complete hash map for the day and then do a central merge. If not done the aggregate computation of top n can be (very) wrong. Imagine top 1 of node 1: {m:5,x:4} node 2 {b:6,x:4}. Aggregate of top 1 is {b:6,m:5 }. In reality the result is {x:8, b:6}
    
    The solution seem to be non deterministic algorithms such as hyperLogLog and family.

* Given that AddLog is called on successive logs it might be possible to optimize it further such that we don't lookup the same log chain at each step

* It doesn't seem realistic to work with more than a few days. Given that logs are contiguous the hash of hash trie implementation could easily be replaced with an big array reducing further startup and look up time.

* it might be possible to not work with urls and replace them with a data structure that is basically a string hash and byte range in the file. Given that the file is memory mapped, it can work nicely, but eventually it will be very disk bound.

* date parsing and conversion to components can be much improved given the standard format 

* performance and especially memory consumption might be further improved using a radix tree instead of the temporary hash map that keeps the url counts

