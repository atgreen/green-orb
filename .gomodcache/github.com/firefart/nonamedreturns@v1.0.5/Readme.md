# nonamedreturns

I hate named returns in golang because they are error prone. That's why I wrote this linter. That's all

Tutorial on how to write your own linter:
https://disaev.me/p/writing-useful-go-analysis-linter/

Named errors used in defers are not reported. If you also want to report them set `report-error-in-defer` to true.

# Why are named returns error prone?

1. **Shadowing of Variables**  
  If you have a named return variable that is also used as a local variable within the function, it can lead to shadowing issues. This occurs when the named return variable is unintentionally shadowed by a local variable with the same name, leading to unexpected behavior.

2. **Accidental Changes**  
   Developers might inadvertently modify the value of a named return variable within the function, thinking it only affects the local variable and not the actual return value.

3. **Readability Issues**  
   While named returns can improve readability for method signatures, they can also make the code harder to understand if misused. It may be unclear whether a variable is a local variable or a return variable, especially in larger functions.

4. **Unused Variables**  
  Named returns often result in variables being declared in the function signature that are not used within the function body. This can lead to confusion and may make the code less maintainable.
