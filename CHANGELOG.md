## [v0.4.0] - 2025-11-10

### Added
- **Enhanced AI prompt generation** – Added explicit rules to omit non‑user‑value content, new sections (BREAKING CHANGES, Deprecated, Security), and stricter categorization, improving prompt relevance and conciseness.

## [v0.3.0] - 2025-11-10  

### Added  
- **Switch to `gpt-oss-120b` model** – The CLI now uses the Cerebras **gpt‑oss‑120b** model as the default for AI‑enhanced release‑note generation. This model offers a larger 65 K token context window and produces more coherent, higher‑quality output, improving the overall usefulness of the generated notes.  
- **Model selection flag (future‑proof)** – Internal code has been refactored to make the model name configurable, paving the way for easy upgrades to newer models without code changes.  

### Changed  
- **README updates** – All documentation now references `gpt‑oss‑120b` as the default model, including usage examples and the AI‑enhancement workflow description.  
- **Internal prompt generation** – The prompt builder now injects the new model name into the AI‑prompt metadata, ensuring the LLM receives the correct context about the inference engine being used.  
- **Token‑budget handling** – The context extraction routine still respects an 8 000‑token budget for relevant files, but the higher‑capacity model can now fully utilise the available window, resulting in richer code‑context snippets.  

### Fixed  
- *(No bug fixes in this release)*  

### Documentation  
- **CHANGELOG entry** – Added a detailed entry for version v0.3.0, describing the model switch and its impact.  
- **README enhancements** – Updated the “Features” and “AI‑Enhanced Release Notes” sections to reflect the new default model and its benefits.  

### Statistics  
- **Files changed**: 9  
- **Commits**: 1  
- **Context analyzed**: ~8 000 tokens  

---  

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
