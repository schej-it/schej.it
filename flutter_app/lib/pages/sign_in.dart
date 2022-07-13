import 'package:flutter/material.dart';
import 'package:flutter_app/models/auth_service.dart';
import 'package:provider/provider.dart';

class SignInPage extends StatelessWidget {
  const SignInPage({Key? key, required this.onSignIn}) : super(key: key);
  final Function(bool signedIn) onSignIn;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: ElevatedButton(
          onPressed: () => _signIn(context),
          child: const Text('Sign in'),
        ),
      ),
    );
  }
  
  void _signIn(BuildContext context) {
    final authService = context.read<AuthService>();
    authService.authenticated = true;
    onSignIn.call(true);
  }
}