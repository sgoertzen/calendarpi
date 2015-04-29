var myStepDefinitionsWrapper = function () {


  this.Given(/^the web server has just started$/, function (callback) {

    callback();
  });

  this.When(/^I view the website$/, function (callback) {
    //this.browser.url("https://localhost").call(callback);
    callback();
  });

  this.Then(/^it will ask for an encryption password$/, function (callback) {
    var title = "hello"
    //this.browser.getElement().and.notify(callback);
    callback();
  });
}
module.exports = myStepDefinitionsWrapper;
