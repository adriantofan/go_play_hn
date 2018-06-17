## Performance


### String vs String Pointer

Just for fun, I have tested a version that has a sting pointer instead of a string. I know that string are value types
but, given that strings are imutable they keep a refference to a byte array where is the text. Still I wanted to test 
my assumption. Here are the results:

#### getDistingQueries on the entire db

    String pointer:
    573697       3	 483784677 ns/op	61922493 B/op	   19113 allocs/op
    573697       3	 486602819 ns/op	61933309 B/op	   19165 allocs/op  6.245s
    573697       2	 534955281 ns/op	61948908 B/op	   19240 allocs/op  3.757s

    String
    573697       3	 435576216 ns/op	61925890 B/op	   19130 allocs/op  5.992
    573697       3	 412219861 ns/op	61946690 B/op	   19230 allocs/op  5.831
    573697       3	 444749906 ns/op	61934834 B/op	   19173 allocs/op  

Non pointer version is about 9% faster, without any impact what soever on the memmory usage.
It consumes quite a bit of memmory 61 MB for a full database query count.

#### BenchmarkReadData

    # String pointer
    BenchmarkReadData-8   	       1	1118722055 ns/op	358121792 B/op	 4870361 allocs/op
    BenchmarkReadData-8   	       1	1119643516 ns/op	358123712 B/op	 4870363 allocs/op 1.145
    BenchmarkReadData-8   	       1	1167624803 ns/op	358123648 B/op	 4870362 allocs/op 1.196
    BenchmarkReadData-8   	       1	1195088345 ns/op	358125696 B/op	 4870366 allocs/op 1.222

    #String
    BenchmarkReadData-8   	       1	1121361311 ns/op	462216960 B/op	 4870364 allocs/op 1.145s
    BenchmarkReadData-8   	       1	1174937320 ns/op	462216768 B/op	 4870361 allocs/op 1.200
    BenchmarkReadData-8   	       1	1181765164 ns/op	462218816 B/op	 4870365 allocs/op 1.209
    BenchmarkReadData-8   	       1	1208232719 ns/op	462218880 B/op	 4870366 allocs/op 1.235

The speed seems to be very marginaly slower on the string version. On the other hand, it consumes quite a bit of memmory:
358 MB versus 462 MB . So the string version has a overhead of arround 25%. 

Conclusions: I will take the speed increase with a slightly more memmory usage(which is anyhow quite high).

Distinct requests are quite expensive. Without changing the algorithm, given that once a result is computed it doesen't 
change I would definitevely use caching. Just by putting a caching proxy in front and setting a far fututre cache expiry 
date on each response should help a lot.