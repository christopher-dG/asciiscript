# asciiscript

Create [asciicasts](https://asciinema.org) without your fingers getting in the way.

Ever tried to record the perfect demo, but couldn't stop missing keys and having to restart?
`asciiscript` lets you record pre-scripted terminal sessions that look human.

## Example

First, create a script.

```sh
echo "Hello, world..."
echo "Here's a demo of asciiscript."

# Comments with a '$' are control commands.

#$ delay 100  - Time between keypresses for subsequent commands (milliseconds).
echo "We can type slow..."
#$ delay 10
echo "Or quite fast."

#$ wait 100  - Time between commands for subsequent commands (milliseconds)
sleep 1  && echo "We can wait for output..."
#$ wait 10
echo "Because otherwise, things could get a bit weird."

echo "I hope you like it!"
```

Then run it with `asciinema`.

```sh
$ asciiscript demo.sh demo.cast
```

[![asciicast](https://asciinema.org/a/207980.png)](https://asciinema.org/a/207980)

## Arguments

`asciiscript` passes on all arguments it receives to `asciinema rec`, except for the script filename.

### Note

There's currently a bug that stops the upload prompt from showing until after you respond to it.
Therefore, it's recommended to use `asciiscript` with a file argument, like so:

```sh
$ asciiscript script.sh output.cast
$ asciinema upload output.cast
```
