# pinched

![image](https://cdn.halceon.dev/img/pinched.jpg)

This is my first ever Go project, so this most likely will take a while. But at least it'll work, I hope?

Pinched aims to be a simple implementation of an "anything" agent to do stuff I need to do. That means:

- Support for Discord and Slack (the two of which I'm most familiar with)
- Support for making PRs to GitHub repos for mini fixes or reading them to examine what's in them
  - considering doing this through straight up claude code automation
- Completions through HCAI (https://ai.hackclub.com)
- Some sort of dashboard with HTMX
- Internal system for how it's called, so it tells itself when to run, but not sure how it'll work yet

The idea started from OpenClaw, but after seeing how incredibly expensive and overcomplicated it felt to set up, I wanted to make something simpler for my own usecase. Not sure how OpenClaw works, but I plan to have it manage its own invocations, which should makei it faster.
