package server

const index = `
<!DOCTYPE html>
<!--[if lt IE 7]>      <html class="no-js lt-ie9 lt-ie8 lt-ie7"> <![endif]-->
<!--[if IE 7]>         <html class="no-js lt-ie9 lt-ie8"> <![endif]-->
<!--[if IE 8]>         <html class="no-js lt-ie9"> <![endif]-->
<!--[if gt IE 8]><!--> <html class="no-js"> <!--<![endif]-->
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
        <title>Shoehorn</title>
        <meta name="description" content="">
        <meta name="viewport" content="width=device-width">

        <link rel="stylesheet" href="/css/application.css">

        <script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.0.8/angular.min.js"></script>
        <script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.0.8/angular-resource.min.js"></script>
        <script src="/js/application.js"></script>
    </head>
    <body ng-app='Shoehorn'>
        <form action='/clone' method='post' >
          <label for=repo>Git Repo</label>
          <input type='text' name='repo'/>
          <input type='submit' value='do it'/>
        </form>
        <div ng-controller="ListCtrl" >
          <div id=site-list >
            <h1> Apps </h1>
            <ul ng-repeat="site in sites">
              <li ng-click="showSite(site)">
              <h3>{{site.name}}</h3>
              </li>
            </ul>
            <h3 ng-show="commandResult">Output</h3>
            <div>
              <p ng-bind-html-unsafe="commandResult.output" ></p>
            </div>
          </div>

          <div id="site-view-container" ng-show="selectedSite">
            <h2>{{selectedSite.name}}</h2>
            <div ng-repeat="process in processes">
              <ul ng-repeat="(name,settings) in process">
                <li>
                  <h3>{{name}}</h3>
                  <div>
                    <b>Actions</b>
                    <button ng-click=start(name)>Start</button>
                    <button ng-click=stop(name)>Stop</button>
                    <button ng-click=restart(name)>Restart</button>
                    <button ng-click=kill(name)>Kill</button>
                  </div>
                  <div>
                    Show Settings: <input type="checkbox" ng-model="checked"/>
                  </div>
                  <ul ng-repeat="(key,val) in settings" ng-show="checked">
                    <li>
                      <i>{{key}}: {{val}}</i>
                    </li>
                  </ul>
                </li>
              </ul>
            </div>
          </div>
        </div>

    </body>
</html>
`
