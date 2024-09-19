# Contributing to Shawty

Thank you for considering contributing to Shawty! We welcome contributions from the community and are excited to see what you'll bring to the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)
- [Submitting Pull Requests](#submitting-pull-requests)
- [Coding Standards](#coding-standard)

## Getting Started

To get started with contributing to Shawty, follow these steps:

1. **Fork the repository**: Click the "Fork" button at the top right of the repository page to create a copy of the repository in your GitHub account.
2. **Clone the repository**: Clone your forked repository to your local machine using the following command:
   ```sh
   git clone https://github.com/wavly/shawty.git
   ```
3. **Create a new branch**: Create a new branch for your changes:
   ```sh
   git checkout -b my-feature-branch
   ```
4. **Install dependencies**: Navigate to the project directory and install the necessary dependencies:
   ```sh
   cd shawty
   go mod tidy
   ```
5. **Make your changes**: Implement your changes in the codebase.
<6. **Run tests**: not available yet

## Reporting Bugs

If you find a bug in Shawty, please report it by creating a new issue in the GitHub repository. Include the following information:

- A clear and descriptive title.
- A detailed description of the bug.
- Steps to reproduce the bug.
- Expected behavior.
- Actual behavior.
- Any relevant screenshots or logs.

## Suggesting Features

We welcome feature suggestions! To suggest a new feature, please create a new issue in the GitHub repository and include the following information:

- A clear and descriptive title.
- A detailed description of the feature.
- The problem the feature would solve.
- Any relevant examples or use cases.

## Submitting Pull Requests

To submit a pull request, follow these steps:

1. **Commit your changes**: Commit your changes to your branch with a clear and descriptive commit message:
   ```sh
   git add .
   git commit -m "update: Add my new feature"
   ```
   WE USE [CONVENTIAL COMMITS STANDARD](https://www.conventionalcommits.org/en/v1.0.0/)
2. **Push your changes**: Push your changes to your forked repository:
   ```sh
   git push origin my-feature-branch
   ```
3. **Create a pull request**: Go to the original repository and click the "New Pull Request" button. Select your branch and provide a clear and descriptive title and description for your pull request.

## Coding Standard
- To maintain consistency please use practical function and variable names
- Write comments on large snippets of code to explain to any future readers
- Format your code using gofmt