<html>
<head>
  <base href="{{.Prefix}}"/>
  <link rel="stylesheet" href="{{.Prefix}}display.css" type="text/css" />
  <script type="text/javascript" src="{{.Prefix}}jquery-3.2.1.min.js"></script>
  <script type="text/javascript" src="{{.Prefix}}display.js"></script>
  <script type="text/javascript">var panel={{.Data.Panel}};</script>
</head>
<body>
  <form>
    <table>
      <tr>
        <td>Universe</td><td><input name="universe" type="number" /></td>
      </tr>
    </table>
  </form>
</body>
</html>
