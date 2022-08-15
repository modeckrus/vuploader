## 0.0.4+100

* -X GET FIX
## 0.0.3+100

* Add `convertFormData` param (defaults: `true`) to convert `FormData` to plain `Map` so we can get
  a CURL representation even while using `FormData` (as for file uploads).

## 0.0.2+100

* Use `dart:developer.log()` instead of `print` for cleaner logs (on android mostly)
* Fix `iOS` example build.

## 0.0.1+100

* Initial release including basic curl requests logging
