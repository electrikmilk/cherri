# Contributing

Before making any contributions, familiarize yourself with how the Cherri programming language generally works for issues, and how the compiler works on a deeper level for PRs.

[Read documentation](https://cherrilang.org/compiler/)

## Issue policy

Before submitting an issue please confirm the following:

- You can reproduce the issue consistently.
- You have searched for your issue.
- You are using the latest version.
- You are describing your issue to the absolute best of your ability. Include any and all relevant information.
- Your issue is one-to-one. No multiple bug issues.
- You are including version info:
  - **If you compiled the compiler:** You are including the version of the Go programming language compiler
  installed on your machine. Use the output of `go version`.
  - **If you are using the latest release:** You are including the version of macOS installed on your machine and your CPU architecture. Use the output
  of `uname -a`.

Confirming these things helps immensely to prevent duplicates, and ensure proper communication and reproducibility.

Issues can be frustrating, but an issue is not easily fixed unless a concrete description of what the issue is, how to
reproduce it and what environment it may be isolated to are clearly communicated.

## Pull request policy

### Programming language policy

- Submit code written in the latest version of the Go programming language that successfully complies using the official
  Go compiler.

### Code style policy

- Use camelCase as is standard for Go code for variable, function, type, etc. names.
- Only add dependencies that add features that exceed the complexity of the compiler, otherwise, you could probably
  write that code yourself.
- This compiler is meant to be one static distributed binary. Do not make it depend on external resources.
- Use the syntax `var variable = ...` rather than `variable := ...`, unless in a for loop or if statement.
- Omit types when they can be inferred, but declare types with no value. (e.g. `var text string`, not: `var text = ""`)
- Use capitalized case for tokens (e.g. `LeftBrace`).
- Only put comments in code to explain _why_ something does something, not what it is doing. Your code should be
  readable enough that it is obvious to anyone reading it.
- No moving code around without a good reason. Only organize code when absolutely necessary.
- Keep any additions to Cherri syntax consistent with the existing style (e.g. camelCase, scripting-language-like, only
  capitalize globals, minimal parenthesis, concise keywords).
- Set the default for action parameter to the default in Shortcuts or best option for the user's privacy. For example, if it's a record audio action, default to recording when the user taps the record button.
- In action names, `detail` is preferred over `component`, and `make` is preferred over `generate` or `create`.

### Code submission policy

- Please use the [Go linter](https://golangci-lint.run/). You can set this up in [GoLand](https://plugins.jetbrains.com/plugin/12496-go-linter). Don't lose focus on your change, but ensure your code is optimal.
- Test your code to the best of your ability, do not submit code that does not compile.
- Write commit titles in the imperative (Fix bug, Add thing, etc.)
- Don't over-explain but don't be vague, briefly but clearly state your changes.
- If you change how a feature works, modify the existing feature, or add a new feature, add or update a related file in the tests folder to add to the commit checks.
- Test your feature or bug fix by writing a Cherri file and checking that it compiles to a valid Shortcut. If you are on
  a non-macOS platform, the Shortcut will compile but not try to sign, so the format may be invalid. Use an XML linter or validator on the debug plist file you can generate using `--debug`.
- Squash commits when doing fixups so that if you remove something, you don't have a commit where you created it
  and then another commit where you remove it, and squash them together so that whatever ended up not being needed is also
  removed from the commit history.
- Make sure your fork is up-to-date every time you are about to submit a contribution. It's best to not make a branch for
  your feature until you are ready to submit it so that you are able to sync your fork beforehand. If this ends up
  happening just sync your main branch with upstream and rebase your feature branch from your main branch.
- No merge commits.
- Commits and changes should be one-to-one, as in, every commit should correlate to a major change in your contribution.
  However, don't make commits for every little thing. Most minor pull requests can contain only one commit.
- Wrap commit messages at 72 characters.
- No periods at the end of commit message titles
- Only capitalize the first word of a commit message title.
- Optionally, add your copyright to the top of files you have changed or added in the same style as the rest of the project.

## Cross-platform development policy

#### It is not encouraged to make major contributions to this project if you do not have access to a Mac.

Troubleshooting why your contribution does not work will be cumbersome. If you prefer to develop on another platform,
that's fine, but it is preferred that you have access to a Mac. For Shortcut debugging, you must use a Mac to validate
that Shortcuts still successfully sign with your changes applied.

This is mainly because although the generated plist may be valid XML, it may be an invalid Shortcut format, which the
`shortcuts` binary on Mac validates when it signs the Shortcut.

If you do not have access to a Mac, you must still make sure that a Cherri file that demonstrates your feature or bug
fix successfully compiles to a valid plist. Use an XML linter or validator to check this and a member of the project
upon reviewing your PR will check if the Shortcut signs successfully.

Please thoroughly debug the issues with your contribution, not repeatedly asking the reviewer to check if a Shortcut signs
with your latest change.

#### Making major contributions without access to a Mac _will_ be frustrating.

---

<p align=center>
  <b>Your issues or PRs will be addressed as time allows for members of this project</b>
</p>

---

## Human language policy

- Keep discussion in issues and pull requests on the topic that the issue or PR addresses.
- Write your comments, pull requests, issues, commits, and code in standard technical English and maintain a professional
  technical tone to the best of your ability. Tools that aid this are recommended to help make this easier. Some IDEs
  even have spell checkers, please don't ignore them.
- Explain things clearly using plain English to the best of your ability.
- Be honest when you don't know something and ask questions, no one is expecting you to "sound
  smart" or pretend to know everything.
- If for some reason you need to specify a unit of measurement please add the alternative as well (e.g. in/cm).

## Burnout

This isn't a requirement, but a suggestion. If you've found that you've thoroughly explored a way to contribute to the
project, try another area, that way you won't burn yourself out doing the same thing. It is also understood and
encouraged to take breaks from the project from time to time.

## No political, religious or other ideological changes

This is a technical project only. We do not make changes or decisions based on our personal religious, political,
or other ideological beliefs. We make changes based on technology and technical standards. Any contributions or issues
that appear to be of this nature will be rejected.
