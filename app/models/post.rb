class Post < ApplicationRecord
  belongs_to :user

  validates :title, presence: true, length: { in: 10..50 }
  validates :content, presence: true, length: { minimum: 20 }
end
