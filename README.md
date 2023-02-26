# golibsignal
![build status](https://github.com/jckimble/golibsignal/actions/workflows/build.yml/badge.svg?branch=master)

golibsignal is a library for accessing signal servers. Written cause I wasn't happy with the state of the textsecure library.

---
* [Install](#install)
* [Roadmap](#roadmap)
* [License](#license)

---

## Install
```sh
go get -u github.com/jckimble/golibsignal
```

## Roadmap
list of features I intend to add over time. volunteers appreciated.
 * TypingMessage - Allow receiving and sending of the TypingMessage
 * CallMessage - Allow receiving calls along with starting them.
 * Device Provisioning - Allow connecting to an existing account.
 * Rewrite axolotl - Get rid of un-needed stores.
 * Profiles - Allow setting profile information from library.

## License

Copyright 2019 James Kimble

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
