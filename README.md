# Go Utils

Common utilities/function used within many of the ContainX projects

## Packages

### encoding

Provides easily exchangeable JSON and YAML based encoding and decoding.

Objects declared using the `json` meta value are compatible with YAML encoding/decoding.

### httpclient

Easy HTTP wrapper which offers simple encoding/decoding using the `encoding` package

### logger

Extends `logrus` offering category based loggers to allow the multi-module
projects to have different logger levels per category

### mockrest

Provides an HTTP server which serves custom responses or status codes.  We use this for testing REST clients

### envsubst

Replaces `${TOKEN}` values from a io.reader into a value matching the key.  

Provides OOB resolver to replace tokens with environment variables. Provides the ability
to specify a custom resolver to in-place lookups when `${TOKENS}` are found

## License

This software is licensed under the Apache 2 license, quoted below.

Copyright 2016 ContainX / Jeremy Unruh

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations under
the License.
