import 'package:flutter/material.dart';
import 'package:flutter_app/models/auth_service.dart';
import 'package:provider/provider.dart';

class SignInPage extends StatefulWidget {
  final VoidCallback onSignIn;

  const SignInPage({
    Key? key,
    required this.onSignIn,
  }) : super(key: key);

  @override
  State<SignInPage> createState() => _SignInPageState();
}

class _SignInPageState extends State<SignInPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: ElevatedButton(
          onPressed: _signIn,
          child: const Text('Sign in'),
        ),
      ),
    );
  }

  void _signIn() async {
    final authService = context.read<AuthService>();
    bool success = await authService.signIn();
    if (success) {
      widget.onSignIn();
    } else {}
  }
}
