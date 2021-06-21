# domz-go-iter1

This is an entry for SCE Programming Challenge #1, put forward by dparrish@. It is written
in Go.

## Usage

`./domz-go-iter1 [BUCKET]`

## Design

A crawling traverse of each synthetic "directory" tree of the bucket is used to
efficiently examine the objects. As directories are discovered, they are listed. New
directories discovered within are, of course, listed, ad infinitum. New objects are checked
against the target ACL, and if they are found to have it, the object name is printed.

The crawling algorithm should allow concurrent listing of any bucket of practical directory
shape. A balanced directory structure will be more efficient -- ideally one with a branching
factor that's a nice multiple of the client machine's core count and a height of 1, 
resulting in even shards that should minimize list overhead. Worst case is lots of
directories with single files in them.

## Runtime

On a standard 8-core "Cloudtop" GCE instance, runtime was approximately 2m13s over three 
test runs. This scanned 214367 files with 14 positive results.
