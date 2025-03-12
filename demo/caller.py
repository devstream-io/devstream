import os

# Set CHATFLOW_BASE to /demo/base
demo_base = os.path.join(os.path.dirname(__file__), "base")
os.environ["CHATFLOW_BASE"] = demo_base

import chatflow  # noqa: E402


def main():
    # res = chatflow.call("/hello_world", "--help")
    # print(res)
    _ = chatflow.call("/hello_world", "This is some user input.")


if __name__ == "__main__":
    main()
