# Wikipedia-Game (Solver)

## What is the Wikipedia Game?
"The Wikipedia game, also known as the Wiki Race or the Wikipedia Challenge, is an online game where players start on a specific Wikipedia page and try to get to another specified page by only clicking on the hyperlinks within Wikipedia articles. The objective is to get to the target page in the fewest number of clicks possible." - ChatGPT (edited)

## 🎯 My objective
The goal of this project is to find the path with fewer clicks from one page to another.

> Basically run a BFS on wikipedia.

## 📥Database Cache
I've created a SQL Server database to save the wikipedia pages that the program visits. It prevents the program from doing many requests to the wikipedia site.

(I might use this wikipedia graph in another project)

## 🔝 Adjustments and improvements

The project already works, but there are a few improvments I'm working on:

- Run requests of the same BFS level in concurrency


## 💻 Skills used

* Requests in Go
* Concurrency in Go
* SQL Server database maintenance
* Regular Expression


<!-- 
[⬆ Voltar ao topo](#wikipedia-game-solver)<br> 
-->
