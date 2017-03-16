# pc - Parser Combinator in go

A parser can be constructed by parser, combinator, and transformer.

* **Parser** parses some text that is read from `Reader`.
* **Combinator** generates new parser from parsers.
* **Transformer** is a kind of combinator but it only does transform the value of the parsed result to any value.
