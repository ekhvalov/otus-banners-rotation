# Banners rotation service
A final project of the **Golang Professional** [course](https://otus.ru/lessons/golang-professional).

The **Banners rotation** service is created for select banners with the highest click-through ratio when the banners 
list and users' preferences change dynamically.

The system is based on the **UCB1** algorithm which is created to solve *Multi-Armed Bandit Problems*.

## Entities

### Slot
A *Slot* is a place on a website where a *Banner* can be shown.
- ID
- Description

### Banner
A *Banner* is an advertisement element that is selected to be shown for a *Social group*.

- ID
- Description

### Social group
A *Social group* is a segmented group of website visitors. For example "women 20-25 y.o."

- ID
- Description
