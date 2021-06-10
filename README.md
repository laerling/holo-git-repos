Entity files are of the format:
```
url=https://github.com/some-user/some-repo
path=/home/user/some-repository
revision=master
```

Note that in place of master you can specify an aritrary revision.
This way you can guarantee what the the repository contents are.


# TODO
- Use logging instead of printf debugging
  - especially for names of temporary files
- delete temporary files at the end of each test
- package
  - executable as /usr/lib/holo/holo-git-repos
  - /etc/holorc.d/git-repos
