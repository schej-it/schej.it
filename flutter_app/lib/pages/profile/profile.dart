import 'package:flutter/material.dart';
import 'package:flutter_app/models/auth_service.dart';
import 'package:provider/provider.dart';

class ProfilePage extends StatelessWidget {
  const ProfilePage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Center(
      child: ElevatedButton(
        onPressed: () {
          final authService = context.read<AuthService>();
          authService.authenticated = false;
        }, 
        child: const Text('Sign out'),
      ),
    ); 
  }
}

