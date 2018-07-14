Article API
--

## Example

#### Using `curl`

`curl -i -X GET http://localhost:4321/articles/1`

`curl -X POST http://localhost:4321/articles -d "@article_data.json"`


`curl -X GET http://localhost:4321/tags/abc/20180712`

{"tag":"abc","count":0,"articles":[],"related_tags":[]}
