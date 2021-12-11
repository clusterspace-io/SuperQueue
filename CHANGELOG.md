### 0.2.2

- Removed exponential back off support since this can be calculated by the consumer using the `attempts` field from the GET and the nack delay.

### 0.2.0

- Added optional delay on POST /nack
- Fixed lock bug
- Added attempts to get record response

### 0.1.0

- Added changelog and version file
