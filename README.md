# forum

## description
This is a project about creating a web forum that allows:


- communication between users
- associating categories to posts
- liking and disliking posts and comments
- filtering posts


## port difficulties

When you attempt to run the server e.g

```bash
go run .
```

then you face port difficulties, i.e the port is currently in use, you can switch to a different port easily by creating an environment variable called PORT.

```bash
export PORT=9000
```

you can now safely restart the server and it will use the port you just provided.

```bash
go run .
```

