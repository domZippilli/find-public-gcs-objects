# find-public-gcs-objects

This program finds objects in a GCS bucket with ACLs that allow public access to the object.

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

On an n1-highcpu-16 GCE instance, runtime was approximately 1m20s over three 
test runs. This scanned 214367 files with 14 positive results.
