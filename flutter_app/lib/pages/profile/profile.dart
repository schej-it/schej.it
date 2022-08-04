import 'package:flutter/material.dart';
import 'package:flutter_app/models/auth_service.dart';
import 'package:provider/provider.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({Key? key}) : super(key: key);

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        _buildProfileSection(),
        Center(
          child: ElevatedButton(
            onPressed: () {
              final authService = context.read<AuthService>();
              authService.signOut();
            },
            child: const Text('Sign out'),
          ),
        ),
      ],
    );
  }

  Widget _buildProfileSection() {
    return Center(
      child: Column(children: [const Text("Tony Xin")]),
    );
  }
}
