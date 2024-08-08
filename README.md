# Paper-Tracker

Paper-Tracker 是一款精简实用程序，旨在监控 arXiv.org 上计算机科学类别的最新研究论文，并提供其基本细节的中文翻译。

## 功能

1. 从 arxiv.org 跟踪 cs.CR 中的最新论文，检索论文标题、发表日期、更新日期、PDF 链接和摘要。
2. 将论文标题和摘要翻译成中文。
3. 在数据库中记录原文和译文信息。
4. 在终端上打印中文版论文信息。

## 技术

1. 使用 GoLang 开发，可与 arxiv.org API 接口。
2. 本地使用 Llama 3.1 进行英汉翻译，无需依赖在线 API。
3. 使用 SQLite 管理记录，以简化数据管理的复杂性。

## 示例

![sample](./output_sample.png 'sample')

---

Paper-Tracker is a streamlined utility designed to monitor the latest research papers from the Computer Science category on arXiv.org and to provide a Chinese translation of their essential details.

## Features

1. Tracks recent papers in cs.CR from arxiv.org, retrieving paper title, publication date, update date, PDF link, and summary.
2. Translates paper titles and summaries into Chinese.
3. Records original and translated information in a database.
4. Prints information of Chinese-version papers on the terminal.

## Techniques

1. Developed in GoLang to interface with the arxiv.org API.
2. Utilizes Llama 3.1 locally for English-to-Chinese translation without relying on online APIs.
3. Manages records using SQLite to simplify data management complexities.

## Sample

![sample](./output_sample.png 'sample')