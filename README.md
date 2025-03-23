# fixfiles

fixfiles is a Go application that helps with debugging by gathering file dependencies for a given error file. It analyzes import statements and includes all relevant files in a single output for easier sharing with LLMs or fellow developers.

## Features

- Recursively analyzes imports to find all dependencies
- Supports multiple programming languages and frameworks:
  - JavaScript/TypeScript (including React, Vue, etc.)
  - Python
  - Go
  - HTML/CSS
  - PHP, Ruby, Java, and more
- Handles different import styles:
  - Relative imports (`./components/Button`)
  - Absolute imports (`/src/utils`)
  - Aliased imports (`@/lib/api`)
- Creates a well-formatted output file with all dependencies

## Installation

1. Ensure you have Go installed on your system (v1.24 or newer)
2. Ensure your `PATH` includes `/usr/local/bin`
2. Clone this repository
3. Build and install the application:

```bash
cd fixfiles
make install
```

## Usage

The basic usage is:

```bash
fixfiles PATH
```

```bash
# Copy output to clipboard
fixfiles PATH | pbcopy
```

```bash
# Save output to file
fixfiles PATH > output.txt
```

Where `PATH` is the path to the file you're having an error with.

### Example

Let's say you're getting an error in `src/components/ClimateInsightsModal.tsx`. Run:

```bash
fixfiles src/components/ClimateInsightsModal.tsx
```

This will:
1. Analyze `ClimateInsightsModal.tsx` for imports
2. Find and include files like `AuthContext.tsx` and `climate.ts` that it depends on

The output will look like:

```
{{ BEGIN CONTENTS OF /path/to/src/components/ClimateInsightsModal.tsx }}
// File contents here
{{ END CONTENTS OF /path/to/src/components/ClimateInsightsModal.tsx }}

{{ BEGIN CONTENTS OF /path/to/src/context/AuthContext.tsx }}
// File contents here
{{ END CONTENTS OF /path/to/src/context/AuthContext.tsx }}

{{ BEGIN CONTENTS OF /path/to/src/lib/api/climate.ts }}
// File contents here
{{ END CONTENTS OF /path/to/src/lib/api/climate.ts }}

------------------------------
Generated 156 lines of code from 3 files
```

## Limitations

- External dependencies (like node_modules) are not included
- Some complex import patterns might not be detected
- Circular dependencies are handled (each file is only included once)

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.
