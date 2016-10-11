# karmahub

Compares the amount of issues and pull requests you created with the amount
of comments and code reviews you did.

The idea is to use it at your daily job organization, so you can get an idea
of how much are you actually contributing to the code review practices.

For example, at the company I work, the "rule of thumb" is
"for each pull request you open, do 3 code reviews".

I can check my progress like this:

```console
./karmahub --user caarlos0 --filter "user:ContaAzul is:pr"
Action    	1m	2m	3m
Authored	78	87	39
Reviewed	119	178	104
```

So, in this scenario, its clear that I'm not following the rule.

Hope this helps you and your team improve the code review practices!
