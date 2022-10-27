package test

const OkTestData = `
<!DOCTYPE html>
<html lang="en" id="current-style">
  <body class="css-ki1d5z">
	<div class="css-4gdbui">
    <div class="css-dxz0om">                  
      <div class="css-tel74u">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-c3r4j7">Author 1</a>
          </div>                      
          <div class="css-1n08q4e">
            <a href class="css-1u6dh35">Oct 9, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-16fbhyp">
        <h1 class="css-1j1qyv3">
          <a href="Link 1" class="css-4zleql">Title 1</a>
        </h1>
        <p class="css-1072ocs">
          <a href class="css-4zleql">Text 1…</a>
        </p>
      </div>
    </div>
  </div>
  <div class="css-4gdbui">
    <div class="css-dxz0om">                  
      <div class="css-tel74u">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-c3r4j7">Author 2</a>
          </div>                      
          <div class="css-1n08q4e">
            <a href class="css-1u6dh35">Sep 8, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-16fbhyp">
        <h1 class="css-1j1qyv3">
          <a href="Link 2" class="css-4zleql">Title 2</a>
        </h1>
        <p class="css-1072ocs">
          <a href class="css-4zleql">Text 2…</a>
        </p>
      </div>
    </div>
  </div>
  <div class="css-4gdbui">
    <div class="css-dxz0om">                  
      <div class="css-tel74u">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-c3r4j7">Author 3</a>
          </div>                      
          <div class="css-1n08q4e">
            <a href class="css-1u6dh35"></a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-16fbhyp">
        <h1 class="css-1j1qyv3">
          <a href="Link 3" class="css-4zleql">Title 3</a>
        </h1>
        <p class="css-1072ocs">
          <a href class="css-4zleql"></a>
        </p>
      </div>
    </div>
  </div>
  <div class="css-4gdbui">
    <div class="css-dxz0om">                  
      <div class="css-tel74u">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-c3r4j7">Author 4</a>
          </div>                      
          <div class="css-1n08q4e">
            <a href class="css-1u6dh35">Oct 20, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-16fbhyp">
        <h1 class="css-1j1qyv3">
          <a href class="css-4zleql">Title 4</a>
        </h1>
        <p class="css-1072ocs">
          <a href class="css-4zleql">Text 4…</a>
        </p>
      </div>
    </div>
  </div>
</div>   
</body>
</html>`

const NoArticlesTestData = `
<!DOCTYPE html>
<html lang="en" id="current-style">
  <body class="lalala"></body>
</html>`

const NoFieldsTestData = `
<!DOCTYPE html>
<html lang="en" id="current-style">
  <body class="css-ki1d5z">
	<div class="css-4gdbui"></div>   
</body>
</html>`

const NoCorrectArticlesTestData = `
<!DOCTYPE html>
<html lang="en" id="current-style">
  <body class="css-ki1d5z">
	<div class="css-4gdbui">
    <div class="css-dxz0om">                  
      <div class="css-tel74u">
        <div class="css-2wkyxu">
          <div class="css-1ajtyzd">
            <a href class="css-c3r4j7">Author 1</a>
          </div>                      
          <div class="css-1n08q4e">
            <a href class="css-1u6dh35">Oct 9, 2022</a>
          </div>
        </div>
      </div>
    </div>
    <div class="css-1wg9be8">
      <div class="css-16fbhyp">
        <h1 class="css-1j1qyv3">
          <a href class="css-4zleql">Title 1</a>
        </h1>
        <p class="css-1072ocs">
          <a href class="css-4zleql">Text 1…</a>
        </p>
      </div>
    </div>
  </div>
  </body>
</html>`
