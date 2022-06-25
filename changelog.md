# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

### Added :
> for new features.

### Changed  :
> for changes in existing functionality.

### Deprecated :
> for soon-to-be removed features.
### Removed :
> for now removed features.
### Fixed :
> for any bug fixe
### Security :
> in case of vulnerabilities.

## [0.1.0]

### Added :
- Blockchain basical mechanics #1 
- Proof of Work algorithm implemented #2 
- CLI that can helps to add some blocks or visualize them #3 
- Blocks persistence #3
## [Unreleased]
## [0.1.1]
### Added :
- Add readme and changelog file
## [0.2.0]
### Added :
- Add transaction system without UTXO persistence layer
- Add new command allowing a user to send tokens into other user
- Add new command allowing to view the tokens balancr of a user
### Changed  :
- Update the way the blockchain was initialized ( init )
### Removed :
- Remove the possibility to add some blocks with the CLI because now a transaction system was implemend so the data arent't accessible with the old way because they are replaced by the transactions inside each block.
## [0.2.1]

Changed : 
- Updating error handling and fixing insufficient balance error bugs