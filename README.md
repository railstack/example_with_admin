# An Example: Integrating Go APIs in a Rails Project

This is a little more advanced example for [go-on-rails](https://github.com/goonr/go-on-rails) generator. There's a [simple one](https://github.com/goonr/example_simple) if you're new to go-on-rails.

The example shows:
* How to use the gems [devise](https://github.com/plataformatec/devise/), [rails_admin](https://github.com/sferik/rails_admin) and [cancancan](https://github.com/CanCanCommunity/cancancan) to implement an admin console to manage a `Post` model
* And use the gem go-on-rails to generate Go APIs codes, make an `index` and a `show` web page

The highlight or primary purpose of the example is using Rails as a tool to `write` Go web app, and give it an admin console in less than 10 minutes.

Now let's go!

## Create a Rails project

At first we create a Rails project:

```bash
rails new example_with_admin --database mysql --skip-bundle
```

Then we edit `Gemfile` to add gems we need:

```ruby
gem 'go-on-rails', '~> 0.1.9'
gem 'devise'
gem 'rails_admin', '~> 1.2'
gem 'cancancan', '~> 2.0'
```

then run `bundle`.

Because we assume you're familiar with Rails, so below we'll not give too much details on Rails part.

Firstly we will follow the `devise` doc to create a user model named `User`, and then we run `rails_admin` and `cancancan` installation also following corresponding steps in their github README doc.

Some points we'll give here is:

1. we add a `role` column to `User` model for `cancancan` to manage the authorization:

```bash
rails g migration add_role_to_users role:string
```

and then edit the migration file to give it a default value: "guest".

2. create a simple `admin?` method in `User` model based the above `role` column:

```ruby
def admin?
  role == "admin"
end
```

## Create a `Post` model

Now let's create a `Post` model:

```bash
rails g model Post title context:text user_id:integer
```

then we edit them as `user` has many `posts` association.

also we can give `Post` model some restrictions, e.g.

```ruby
# app/models/post.rb
class Post < ApplicationRecord
  belongs_to :user

  validates :title, presence: true, length: { in: 10..50 }
  validates :content, presence: true, length: { minimum: 20 }
end
```

## Use rails_admin to create some posts

Then we `rails s` to start the server and visit `http://localhost:3000/admin`, follow the steps to signup a user, but don't login immediately.

We make the user `admin` role in Rails console:

```ruby
User.first.update role: "admin"
```

Then we login the rails_admin to create some posts.

## Generate Go codes and make pages

Now it's time to show the power of `go-on-rails`:

```bash
rails g gor dev
```

Then we go to the generated `go_app` directory to make two web pages, an `index` page to list all posts, and a `show` page to show a post content.

We need make two template files named `index.tmpl` and `show.tmpl` under the `views` directory at first, and create a controller file named `post_controller.go` in controller directory, the controller file have two `handler` functions: a `IndexHandler` and a `ShowHandler`.

We use the `Gin` framework, more details you can refer to its doc: https://github.com/gin-gonic/gin.

And at last edit the `main.go` to add two router paths to pages:

```go
r.GET("/", c.IndexHandler)
r.GET("/posts/:id", c.ShowHandler)
```

Run the server on port 4000:

```bash
go run main.go -port 4000
```

Now you can visit the `index` page on: http://localhost:4000.

## What's next?

We'll do next:
* Add CSS to pages
* Paginate post list in `index` page
* `User` authentication in Go app by JWT tokens...

Welcome back and let's go on ...

## Make a nice page

Now let's make things a little bit complicated, we'll abandon the templates we used above.

Because Rails 5.1 released with a gem [webpacker](https://github.com/rails/webpacker), we'll use it to implement an independent frontend view to call our Go APIs to render an index and a post pages.

We try to make the steps clearer than previous sections due to the necessary details of new techniques.

### Use webpacker to develop an independent frontend

We can install and configure the webpacker by simply following its [README](https://github.com/rails/webpacker) file, but one thing need to mention is to make sure some [prerequisites](https://github.com/rails/webpacker#prerequisites) installed before hand. After that we can begin to install webpacker:

```ruby
# Gemfile
gem 'webpacker', '~> 3.0'
```

then bundle and install:

```bash
bundle
bundle exec rails webpacker:install
```

Because we plan use React in our project, so install the React support of webpacker:

```bash
bundle exec rails webpacker:install:react
```

We can see a some directories and files generated into the Rails, the important directory is the `app/javascript/packs`, we'll write our react programs all there by default webpacker config.

### Install material-ui and react-router-dom

[material-ui](https://github.com/callemall/material-ui) is a package of React components that implement Google's Material Design, use it can make some complex and beautiful pages easily. So we'll try it.

And we have a index and a post pages, so we need to make a frontend router for these different views. We use [react-router-dom](https://github.com/ReactTraining/react-router) to make it.

Now let's install them by the package tool `yarn`(although you can use `npm` as well):

```
yarn add material-ui react-router-dom
```

### We need a Rails view as a server-side template

Now let's create a controller and a page as a server-side template to integrate the webpacker's helper and React files. If you don't understand this currently, just keep it in mind, I'm sure you'll get it after steps.

```bash
rails g controller Pages home
```

Change the route `get 'pages/home'` that Rails added by default in `config/routes.rb` to:

```ruby
root to: 'pages#home'
```

Edit the view `app/views/pages/home.html.erb`, clear all the content of it and add the line below:

```ruby
<%= javascript_pack_tag 'hello_react' %>
```

Open another terminal and change to the our project's root directory, run:

```bash
./bin/webpack-dev-server
```

and meanwhile run `rails s` in the current terminal, then open `http://localhost:3000` in your browser, you can see it prints a "Hello React!". Now we make React and webpacker worked in Rails.

Next let's remove this default installed `hello_react` example from our project:

```bash
rm app/javascript/packs/hello_react.jsx
```

We will use another default installed file `app/javascript/packs/application.jsx` as our main React programs, so we change the `javascript_pack_tag` reference file in our Rails view above from `hello_react` to `application`:

```ruby
<%= javascript_pack_tag 'application' %>
```

then we clear the content of `application.jsx` for further programming. And next we will add two more React files as React components corresponding to our `index` and `show` pages.

### Write React components for `index` and `show` pages

Firstly, let's construct the routes for the two views in `application.jsx`.

```javascript
import React from 'react'
import ReactDOM from 'react-dom'
import { HashRouter, Route, Link } from 'react-router-dom'
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider'
import AppBar from 'material-ui/AppBar';

import IndexCard from './index'
import ShowCard from './show'

// some routes

<Route exact path="/" component={Index}/>
<Route path="/posts/:id" component={Show}/>

```

Here we use `HashRouter` instead of `Router` directive because we need the server-side render at first, `/path` will be an invalid route for Rails while `/#path` will be manipulated by frontend.


And we'll created three React components to get the job done:

* [application.jsx](https://github.com/goonr/example_with_admin/blob/master/app/javascript/packs/application.jsx)
* [index.jsx](https://github.com/goonr/example_with_admin/blob/master/app/javascript/packs/index.jsx)
* [show.jsx](https://github.com/goonr/example_with_admin/blob/master/app/javascript/packs/show.jsx)

you can check the files to get the details.

Here we use a javascript package named `axios` to do Ajax requests, you can install it by yarn:

```bash
yarn add axios
```

### Server side changes

Now when we set up the Rails server to request APIs of our Go application, we need to add a CORS configuration to the Go server to make cross domains accessible. Because we use the Gin framework by default, we choose a cors package specially for the Gin: [github.com/gin-contrib/cors](github.com/gin-contrib/cors).

We just use its default configuration that allows all the Origins access our Go server for testing easily, here's [the details](https://github.com/gin-contrib/cors#default-allows-all-origins).

### Try the new views

In one terminal you set up the Go server in port 4000 under the `go_app` directory:

```bash
go run main.go --port 4000
```

In another terminal run `rails s` to set up the Rails server in default port 3000, and meanwhile to get the `./bin/webpack-dev-server` server up in a third terminal.

Now visit the http://localhost:3000, you can see our new pages.

(WIP)
