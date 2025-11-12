## [v0.7.6] - 2025-11-12

### Added
- **Introducing a Sophisticated 2-Stage Polish Workflow for Premium Release Notes**
  We've implemented an advanced two-stage workflow designed to elevate the quality of your release notes. This new process intelligently combines technical discovery models to accurately identify changes with customer-facing language models to craft clear, engaging, and polished descriptions. This ensures your release notes are both technically precise and easy for your users to understand.

- **Expanded AI Model Access with OpenRouter Support**
  We've integrated support for OpenRouter, significantly expanding your access to a vast array of AI models. This integration allows you to leverage hundreds of different AI models through a single, streamlined API, offering greater flexibility and power in generating your release notes.

- **Configurable File Exclusions for Focused AI Context**
  You can now specify files or directories to be excluded when generating release notes. This enhancement ensures that our AI models focus exclusively on the most relevant code changes, preventing irrelevant information from impacting the accuracy and conciseness of your release notes.

### Changed
- **Updated Default AI Provider to Cerebras (llama-3.3-70b) for Enhanced Performance**
  The default AI provider has been updated to Cerebras, utilizing the powerful `llama-3.3-70b` model. This change provides improved accessibility and delivers higher quality, more accurate results for your release note generation.

- **Granular Configuration for the 2-Stage Polish Workflow**
  The newly introduced 2-stage polish workflow is now highly configurable. This allows you to precisely control its operation and tailor its behavior to meet your specific requirements, giving you more power over the final output.

- **Enhanced GitHub Actions Workflow for Improved Integration**
  Our GitHub Actions workflow has been updated to fully leverage the new `llama-3.3-70b` model and integrate seamlessly with OpenRouter support. This update streamlines your automation processes and ensures you benefit from the latest AI capabilities directly within your CI/CD pipeline.

### Fixed
- **Resolved Context Pollution from Meta-Documentation Changes**
  We've addressed an issue where changes in meta-documentation could inadvertently influence the content of AI-generated changelogs. This fix ensures that your release notes remain focused solely on relevant code changes, preventing unintended information from being included.

- **Accurate Reflection of Significant Code Changes in Release Notes**
  This update resolves an issue to ensure that all impactful code modifications are accurately identified and included in your release notes. You can now be confident that every significant change is properly documented for your users.

- **Correct API Key Environment Variable Overrides**
  We've resolved a problem where overriding the default AI provider via command-line interface (CLI) flags did not correctly update the corresponding API key environment variable. This fix ensures that your specified API keys are always correctly applied, preventing authentication issues.