# Notes
Here are my notes in reviewing/refactoring this code:

1. Idiomatic Go is a thing. Embrace it.

2. All in all, your code was really, really good.  My code review/refactoring is taking you from 80% to 90%, with some of my review/refactors are based on my personal style that AFAIK does not conflict with idiomatic Go but really helps with readability and maintainability.  _(I can't take you to 100% because I'm just not that good; I doubt almost anyone is.)_

3. One of the first things I think developers new to Go do is see packages and immediately start envisioning how their namespaces will allow logical code organization within an applications. And then these developers start creating a lot of packages with generic names. I absolutely did that. Unfortunately however I found through painful outcomes that using packages to organize code within an application is not a smart approach. Generic names often conflict with packages names in the Go standard library, and it is very easy to create cyclical dependencies across your packages which then requires you to work very hard to create internal interfaces you would otherwise never need just to get your over-organized code to compile. In reality your packages are rarely actually independent because one so often cannot really standalone without the functionality of the others. 

4. You had a lot of code in your commands that could really be re-envisioned as methods for your structs which really simplify the code in your commands and allows for reuse of functionality now or down the road. 

5. So for any given application or library package it is better — ignoring the required `main` package — to stick with one, or a very small number of packages, where the `main()` func is in the `main()` package — if you are writing an application — and the rest of code should be in a library named after your application, e.g. `lpm` in your case.  I reorganized your code in my fork to this effect.       

6. Your two main files (`darwin.go` and `windows.go`) were almost exactly the same except for one expression difference. I refactored that difference into `utils` and I made your `main.go` file simple and easy to grok.

7. Your functions should _(almost)_ always return an error in Go when an error occurs that your func cannot handle itself. Just failing to terminal is not a good strategy. You will also be able to add a lot more information so you user can figure out what happened. And if you ever want to use the same packages for microservices you will thank me.

8. So when possible, if you want to fail on an error then bubble up that error while wrapping it with more info then handle at the root, in a CLI. Don't just fail deep in the middle of code because you may later want to reuse that code and if you do it won't be robust and will require lots of rearchitecture.

9. When using text to describe error messages start with lower case so multiple messages can be composed without it looking weird. Also don't use periods to end error messages.

10. Don't create an empty instance to be able to call a method. Just define as a func. That's idiomatic in Go and acceptable. If for want an instance for dependency injection but still want to call it explictly elsewhere, create both a method and a func that creates an instance to be able to call the method. So don't write `util.HTTPUtil{}.Get()`. Just write `util.HttpGet()`

11. I refactored out `else{}` statements. [Embrace the happy-path](https://medium.com/@matryer/line-of-sight-in-code-186dd7cdea88). 

12. It's more idiomatic in Go to use fmt.Sprintf() than string concatenation.
 
13. My personal style is to use `goto end` rather than early `return` because it provides numerous benefits. Unfortunately in Go it also requires all declarations be above the `goto`. But then explict declarations help the reader of your code, so I don't think that's really a bad thing. Also doing it this way causes the developer to have to work harder to write longer methods, so this approach encourages more modularity.
 
14. `errors` is a well-known package name in Go; you probably should not use it for your own package name.
 
15. For cross-platform paths, use `filepath.FromSlash()` and `filepath.ToSlash()`
 
16. I use GoLand from JetBrains as my editor/IDE _(which I cannot praise too much.)_ Thus I added `//goland:noinspection` in a few places because in those places it flagged that it reasonably thought were errors that were not really errors, such as not capturing the return value of a `.Close()` in a `defer` statement.
 
17. In `Dependencies` you use pointer and non-pointer receivers (`d`) in your methods. Better to be consistent and use the same. I know it has less immutability benefit to use pointers, but mixing can cause people who have to use your packages to struggle. If you need immutability code for it explicitly, such as how `append()` works in Go.
 
18. You can define types such as 'ClasspathFiles' so you can type `[]fs.FileInfo` in just one place and so that when you use them only values that are meant to be used for Classpath Files can be used for Classpath files, vs. some other slice of files.
 
19. Your configuration in App variables makes it harder to write unit tests, from experience. Every times I have started with variables for config I ultimately end up having to refactor to using a struct with methods.

20. Things that run within `func init()` should not be of the nature to throw errors. It will make debugging edge cases much harder. Do things that might throw errors in the main flow of your program.

21. It appears that `Dependency` is a `map[string]string` but it looks like it never has more than a single key and a single value. Why not a simple struct? 

22. Test in the same package has too much access to internals. Better to create a separate `test` package

23. In general I have tried to move away from calculating things in init() functions and instead do so on demand in methods, e.g. see `Context.GetManifestFilepath()` vs. `app.FileLocation`. 