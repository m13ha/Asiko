#!/bin/bash

# This script generates a TypeScript Fetch API client from the OpenAPI specification.

# Ensure the generator is installed:
# npm install @openapitools/openapi-generator-cli -g

# Set the input spec file and output directory
INPUT_SPEC="../docs/swagger.json"
OUTPUT_DIR="../api-client/src"

# Remove the old client to ensure a clean build
rm -rf $OUTPUT_DIR

# Generate the new client
openapi-generator-cli generate \
    -i $INPUT_SPEC \
    -g typescript-fetch \
    -o $OUTPUT_DIR \
    --additional-properties=typescriptThreePlus=true,usePromises=true

echo "✅ API client generated successfully in $OUTPUT_DIR"
