# exampledebug
A cli tool to diff output from go test failures for example tests. 

Pipe the output from go test with [example driven tests](https://blog.golang.org/examples) into example debug:

```
$ go test ./pkg/storage/closedts/minprop | exampledebug
Example:
--- got
+++ want
@@ -11,7 +11,7 @@
 All commands initially start out on the right. The command has its timestamp forwarded to 0.000000000,2 .

   closed=0.000000000,0
-      |            next=0.000000000,1
+      |            next=0.000000000,2
       |          left | right
       |             0 # 1
       v               v

---
ExampleTracker_Close:
--- got
+++ want
@@ -12,4 +12,4 @@

 Closed: 2.000000000,0 map[99:1]
 Note how the MLAI has 'regressed' from 2 to 1. The consumer
-needs to track the maximum over all deltas received.
\ No newline at end of file
+needs to track the maximum over all deltas received
\ No newline at end of file

---
```
