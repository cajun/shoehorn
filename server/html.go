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

        <div ng-controller="ListCtrl" >
          <div id=site-list >
            <h1> Apps </h1>
            <ul ng-repeat="site in sites">
              <li ng-click="showSite(site)">
              <h3>{{site.name}}</h3>
              </li>
            </ul>
          </div>

          <div id="site-view-container" ng-show="selectedSite">
            <h2>{{selectedSite.name}}</h2>
            <ul ng-repeat="process in processes">
              <li>
              <i>{{process.App}}</i>
              </li>
            </ul>
          </div>
        </div>

    </body>
</html>
`
