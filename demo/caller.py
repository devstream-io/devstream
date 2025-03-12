import os

# Set DEVSTREAM_BASE to /demo/base
demo_base = os.path.join(os.path.dirname(__file__), "base")
os.environ["DEVSTREAM_BASE"] = demo_base

import devstream  # noqa: E402


def main():
    # res = devstream.call("/hello_world", "--help")
    # print(res)
    _ = devstream.call("/hello_world", "This is some user input.")


if __name__ == "__main__":
    main()
