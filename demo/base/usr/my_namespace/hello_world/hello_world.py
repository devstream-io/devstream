import sys


def main():
    print("Hello world from devstream demo.")
    if len(sys.argv) > 1:
        print(f"User input: {sys.argv[1]}")
    else:
        print("No user input provided.")


if __name__ == "__main__":
    main()
