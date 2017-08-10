# An Example of go-on-rails Generator: Integrating Go APIs in a Rails Project

This is a little more advanced example for [go-on-rails](https://github.com/goonr/go-on-rails) generator. There's a [simple one](https://github.com/goonr/example_simple) if you're new to go-on-rails.

The example shows:
* How to use the gems [devise](https://github.com/plataformatec/devise/), [rails_admin](https://github.com/sferik/rails_admin) and [cancancan](https://github.com/CanCanCommunity/cancancan) to implement a admin console to manage a `Post` model,
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
gem 'go-on-rails', '~> 0.1.4'
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

## User rails_admin to create some posts

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

And at last edit the `main.go` to add two router paths to pages.

## What's next?

We'll do `User` authentication in Go app by JWT tokens...

