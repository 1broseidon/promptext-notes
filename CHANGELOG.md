## [v0.2.0] - 2025-11-10

### Added
- **Automated AI-Enhanced Release Notes Generation**: Added a feature to automatically generate release notes using AI enhancement. This feature utilizes the Cerebras API with the `llama-3.3-70b` model, which has a 65,536 token context window, allowing for more detailed and accurate release notes.
- **Detailed Error Logging**: Added detailed error logging for API responses to improve debugging and error handling.
- **Redirect Status Messages to Stderr**: Added a feature to redirect all status messages to stderr in the release notes script, improving the overall user experience.

### Changed
- **Cerebras API Model**: Changed the Cerebras API model from `llama3.1-8b` to `llama-3.3-70b` to take advantage of the larger context window and improved performance.
- **CI/CD Pipeline**: Made the coverage check non-blocking in the CI/CD pipeline to improve the build process.

### Fixed
- **Bug in Token Counting**: Fixed a bug in token counting to ensure accurate token counting and analysis.
- **Edge Case in File Filtering**: Fixed an edge case in file filtering to improve the accuracy of file analysis.

### Documentation
- **Updated CHANGELOG**: Updated the CHANGELOG to reflect the changes and improvements made in this release.
- **Updated README**: Updated the README to include the latest features, usage, and examples.

### Statistics
- **Files changed**: 9
- **Commits**: 7
- **Context analyzed**: ~7708 tokens

This release focuses on improving the accuracy and detail of release notes, as well as enhancing the overall user experience. The addition of automated AI-enhanced release notes generation and detailed error logging improves the quality and reliability of the release notes. The changes to the Cerebras API model and CI/CD pipeline improve the performance and efficiency of the build process.

🧠 Calling Cerebras API (llama3.1-70b)...

