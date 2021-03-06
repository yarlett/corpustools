#corpustools

corpustools is a Go library which provides statistical natural language processing functionality on a corpus or words / characters / whatever, represented internally as a sequence of integers. It is intended to be fairly lightweight and fast, and to operate on corpora which can be held in RAM (of the order of hundreds of millions, or the low billions, of tokens). With this in mind many of its operations are based on a suffix array index into corpora.

The functionality of the library includes:
	o Fast (O(log(N))) computation of frequency statistics over ngrams in the corpus, to underpin applications such as language modeling.
	o Rapid enumeration of all distinct ngrams in the corpus.
	o Nearest neighbor search over ngrams in the corpus based on their distributional statistics (co-occurrence modeling).
	o Computation of subsequences in the corpus which allow the description length of the corpus to be reduced.

##Installation instructions

1. Install Go on your system.
2. Download this repository into your $GOPATH.
3. Type "go install corpustools".

##Background

Much of the core functionality of the library is available through the Corpus object.

To create a corpus object from a text file:

```go
import "corpustools"
lowerCase, returnChars := true, false
corpus := corpustools.CorpusFromFile("myfile.txt", lowerCase, returnChars)
```

which will tokenize myfile.txt into lower-case words (not characters), assigning each distinct token a unique integer key, and representing the sequence of tokens internally as an array of integers (the Corpus object also encapsulates the mapping from tokens to integers so that the internal representation can be translated back into its original format).

When the Corpus object is created, a suffix array for the corpus is also created. The suffix array consists of integer indexes into the corpus such that traversing the suffix array in order and following its pointers results in an ordered enumeration of all substrings of the corpus. For example, here is an example corpus and its corresponding suffix array:

```go
corpus: [5, 4, 3, 2, 1]
suffix: [4, 3, 2, 1, 0]
```

We can see that suffix[0] points to position 4 in the corpus, which corresponds to the substring [1]. suffix[1] points to position 3 in the corpus, which corresponds to the substring [2, 1] > [1], and so on.

Although creating a suffix array doubles the amount of memory required by the Corpus object, it allows much faster searching of the corpus. A naive search algorithm in order to find the frequency of the substring [1, 2, 3] in the corpus, for example, will take O(N) time, whereas a binary search over the suffix array takes only O(log(N)) time.

##Usage

Once a corpus has been created, it is relatively easy to use. For example, in order to find the locations at which a given sequence is located:

```go
indices := corpus.Find([]int{1, 2, 3})
```

In order to compute the frequency of a specified sequence:

```go
f := corpus.Frequency([]int{1, 2, 3})
```

And to generate the set of all unique, ordered trigrams in the corpus:

```go
trigrams := corpus.Ngrams(3)
```

Further and more detailed examples of the functionality provided by the library are included in the /examples folder.