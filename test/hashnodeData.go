package test

const OK_TEST_DATA = `
<!DOCTYPE html>
<html lang="en" id="current-style">
  <body class="css-ki1d5z">
	<div class="css-1s8wn94">
    <div class="css-dxz0om">                  
      <div class="css-cj4uuj">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-9ssaz8">Author 1</a>
          </div>                      
          <div class="css-1cyn8lj">
            <a href class="css-15gyiyx">Oct 9, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-1abv9a9">
        <h1 class="css-1yrl49b">
          <a href="Link 1" class="css-4zleql">Title 1</a>
        </h1>
        <p class="css-1m4ptby">
          <a href class="css-4zleql">Text 1…</a>
        </p>
      </div>
    </div>
  </div>
  <div class="css-1s8wn94">
    <div class="css-dxz0om">                  
      <div class="css-cj4uuj">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-9ssaz8"></a>
          </div>                      
          <div class="css-1cyn8lj">
            <a href class="css-15gyiyx">Sep 8, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-1abv9a9">
        <h1 class="css-1yrl49b">
          <a href="Link 2" class="css-4zleql">Title 2</a>
        </h1>
        <p class="css-1m4ptby">
          <a href class="css-4zleql">Text 2…</a>
        </p>
      </div>
    </div>
  </div>
  <div class="css-1s8wn94">
    <div class="css-dxz0om">                  
      <div class="css-cj4uuj">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-9ssaz8">Author 3</a>
          </div>                      
          <div class="css-1cyn8lj">
            <a href class="css-15gyiyx">Sep 8, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-1abv9a9">
        <h1 class="css-1yrl49b">
          <a href="Link 3" class="css-4zleql">Title 3</a>
        </h1>
        <p class="css-1m4ptby">
          <a href class="css-4zleql"></a>
        </p>
      </div>
    </div>
  </div> 

</body>
</html>`

const NO_CORRECT_ARTICLES_TEST_DATA = `
<!DOCTYPE html>
<html lang="en" id="current-style">
  <body class="css-ki1d5z">

	<!-- Link is empty -->
	<div class="css-1s8wn94">
    <div class="css-dxz0om">                  
      <div class="css-cj4uuj">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-9ssaz8">Author 1</a>
          </div>                      
          <div class="css-1cyn8lj">
            <a href class="css-15gyiyx">Oct 9, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-1abv9a9">
        <h1 class="css-1yrl49b">
          <a href class="css-4zleql">Title 1</a>
        </h1>
        <p class="css-1m4ptby">
          <a href="" class="css-4zleql">Text 1…</a>
        </p>
      </div>
    </div>
  </div>

	<!-- Link attribute not found -->
	<div class="css-1s8wn94">
    <div class="css-dxz0om">                  
      <div class="css-cj4uuj">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-9ssaz8">Author 1</a>
          </div>                      
          <div class="css-1cyn8lj">
            <a href class="css-15gyiyx">Oct 9, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-1abv9a9">
        <h1 class="css-1yrl49b">
          <a class="css-4zleql">Title 1</a>
        </h1>
        <p class="css-1m4ptby">
          <a href="" class="css-4zleql">Text 1…</a>
        </p>
      </div>
    </div>
  </div>

<!-- Title node not found -->
	<div class="css-1s8wn94">
    <div class="css-dxz0om">                  
      <div class="css-cj4uuj">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-9ssaz8">Author 1</a>
          </div>                      
          <div class="css-1cyn8lj">
            <a href class="css-15gyiyx">Oct 9, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-1abv9a9">
        <h1 class="css-1yrl49b">
          
        </h1>
        <p class="css-1m4ptby">
          <a href="" class="css-4zleql">Text 1…</a>
        </p>
      </div>
    </div>
  </div>

<!-- Title node is empty -->
	<div class="css-1s8wn94">
    <div class="css-dxz0om">                  
      <div class="css-cj4uuj">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-9ssaz8">Author 1</a>
          </div>                      
          <div class="css-1cyn8lj">
            <a href class="css-15gyiyx">Oct 9, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-1abv9a9">
        <h1 class="css-1yrl49b">
          <a href="url" class="css-4zleql"></a>
        </h1>
        <p class="css-1m4ptby">
          <a href="" class="css-4zleql">Text 1…</a>
        </p>
      </div>
    </div>
  </div>
  </body>
</html>`

const ARTICLES_WITH_WARNINGS = `
<!DOCTYPE html>
<html lang="en" id="current-style">
  <body class="css-ki1d5z">

	<!-- Description node not found -->
	<div class="css-1s8wn94">
    <div class="css-dxz0om">                  
      <div class="css-cj4uuj">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-9ssaz8">Author 1</a>
          </div>                      
          <div class="css-1cyn8lj">
            <a href class="css-15gyiyx">Oct 9, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-1abv9a9">
        <h1 class="css-1yrl49b">
          <a href="Link 1" class="css-4zleql">Title 1</a>
        </h1>
        <p class="css-1m4ptby">
          
        </p>
      </div>
    </div>
  </div>

	<!-- Author node not found -->
	<div class="css-1s8wn94">
    <div class="css-dxz0om">                  
      <div class="css-cj4uuj">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            
          </div>                      
          <div class="css-1cyn8lj">
            <a href class="css-15gyiyx">Oct 9, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-1abv9a9">
        <h1 class="css-1yrl49b">
          <a href="Link 2" class="css-4zleql">Title 2</a>
        </h1>
        <p class="css-1m4ptby">
          <a href="" class="css-4zleql">Text 2…</a>
        </p>
      </div>
    </div>
  </div>

<!-- Date node not found -->
	<div class="css-1s8wn94">
    <div class="css-dxz0om">                  
      <div class="css-cj4uuj">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-9ssaz8">Author 3</a>
          </div>                      
          <div class="css-1cyn8lj">
            
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-1abv9a9">
        <h1 class="css-1yrl49b">
          <a href="Link 3" class="css-4zleql">Title 3</a>
        </h1>
        <p class="css-1m4ptby">
          <a href class="css-4zleql">Text 3…</a>
        </p>
      </div>
    </div>
  </div>

  </body>
</html>
`
