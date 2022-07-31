import 'package:auto_route/auto_route.dart';
import 'package:flutter_app/models/auth_service.dart';
import 'package:flutter_app/router/app_router.gr.dart';

class AuthGuard extends AutoRedirectGuard {
  final AuthService _authService;

  AuthGuard(this._authService) {
    _authService.addListener(() {
      if (!_authService.authenticated) {
        reevaluate();
      }
    });
  }

  @override
  void onNavigation(NavigationResolver resolver, StackRouter router) {
    if (_authService.authenticated) {
      return resolver.next();
    } else {
      router.push(SignInPageRoute(
        onSignIn: () {
          resolver.next();
          router.removeLast();
        },
      ));
    }
  }

  @override
  Future<bool> canNavigate(RouteMatch route) async {
    if (!_authService.authenticated) {
      return false;
    }
    return true;
  }
}
