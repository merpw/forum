import { Post, User } from "../custom"

export const user = {
  id: 1,
  name: "Max",
  email: "max@mer.pw",
}

export const users: User[] = [user, { name: "Cat", id: 2 }]

export const posts: Post[] = [
  {
    id: 3,
    title: "Post 3",
    content: "Short one",
    author: users[0],
    date: "2022-12-26T15:58:18.166Z",
    likes: 0,
    dislikes: 0,
    comments: [],
  },
  {
    id: 1,
    title: "Post 1",
    content:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam euismod imperdiet dictum. Sed aliquam lacus neque, in efficitur diam ultrices quis. Nam viverra varius quam, a dictum purus congue eu. Vivamus at vehicula massa. Pellentesque non eros augue. Praesent ut ante quis dui hendrerit iaculis ut pretium velit. Vivamus sed ante aliquam, accumsan ipsum nec, suscipit lorem. Ut eu diam nec mi mattis imperdiet. Vivamus massa est, cursus vitae felis sed, congue semper lectus. Proin placerat gravida nisi, malesuada pulvinar ligula semper non. Sed rhoncus augue id tempus congue. In laoreet a mi nec pharetra. Praesent tempor erat et eros bibendum porttitor at et erat. Proin aliquet pellentesque lacus, non ullamcorper mi placerat pulvinar.",
    author: users[0],
    date: "2022-12-22T19:36:18.166Z",
    likes: 0,
    dislikes: 0,
    comments: [
      {
        author: users[0],
        text: "Donec neque diam, sodales eget aliquam vel, fringilla nec eros. Proin tincidunt felis arcu, a tempus ante hendrerit eget. Vestibulum eget nisi eget tellus porttitor interdum a nec velit.",
        date: "2022-12-22T21:36:18.166Z",
      },
      {
        author: users[1],
        text: "Bruh...",
        date: "2022-12-23T20:36:18.166Z",
      },
    ],
  },
  {
    id: 2,
    title: "Post 2",
    content:
      "Aenean non semper sapien. Sed at tristique sapien. Etiam auctor accumsan mi, a lobortis enim sollicitudin a. Aliquam pellentesque ligula ullamcorper egestas tempor. Curabitur congue nec odio ut ultricies. Donec tempus tincidunt mi non bibendum. Sed quis malesuada ex. Mauris gravida luctus interdum. Vivamus leo ligula, scelerisque sit amet arcu ac, rhoncus rutrum augue. Integer congue magna sem, sit amet euismod lacus blandit sit amet. Maecenas ultrices odio ipsum, eget imperdiet nibh tempus a. Donec neque diam, sodales eget aliquam vel, fringilla nec eros. Proin tincidunt felis arcu, a tempus ante hendrerit eget. Vestibulum eget nisi eget tellus porttitor interdum a nec velit.",
    author: users[1],
    date: "2022-12-22T19:36:18.166Z",
    likes: 0,
    dislikes: 0,
    comments: [],
  },
]
