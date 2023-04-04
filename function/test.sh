#!/bin/bash

curl -XPOST localhost:1234/invoke -d '{"id": "cpu", "args":{"n": , "vectors": [[1,1], [1,1]]}}' -H "Content-Type: application/json"


