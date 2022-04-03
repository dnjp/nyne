# nynetab

Implements tab expansion and indentation.

```
Usage of nynetab:
	nynetab [-unindent]
```

Nynetab is what is used under the hood for tab expansion in nyne.
Executing `nynetab` will insert either a hard or soft tab depending
on [what is configured](https://github.com/dnjp/nyne/blob/master/config.go).
Executing `nynetab -unindent` will unindent text that is selected.
