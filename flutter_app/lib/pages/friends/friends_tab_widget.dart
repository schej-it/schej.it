import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/components/friends/friend_card.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/router/app_router.gr.dart';

class FriendsTabWidget extends StatefulWidget {
  const FriendsTabWidget({Key? key}) : super(key: key);

  @override
  State<FriendsTabWidget> createState() => _FriendsTabWidgetState();
}

class _FriendsTabWidgetState extends State<FriendsTabWidget> {
  // Controllers
  final TextEditingController _searchTextController = TextEditingController();

  // Variables
  var friends = [
    {'name': 'Winston Tilton', 'status': FriendStatus.free},
    {
      'name': 'Samantha Jones',
      'status': FriendStatus.busy,
      'curEventName': 'BTG meeting'
    },
    {
      'name': 'Tyler Smithson',
      'status': FriendStatus.busy,
      'curEventName': 'PSYC 336'
    },
    {'name': 'Winston Tilton', 'status': FriendStatus.invisible},
    {'name': 'Winston Tilton', 'status': FriendStatus.invisible},
    {'name': 'Winston Tilton', 'status': FriendStatus.invisible},
    {'name': 'Winston Tilton', 'status': FriendStatus.invisible},
    {'name': 'Winston Tilton', 'status': FriendStatus.invisible},
  ];

  @override
  void dispose() {
    _searchTextController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.symmetric(
        horizontal: SchejConstants.pagePadding.left,
      ),
      child: Column(
        children: [
          _buildSearchTextField(),
          Expanded(child: _buildFriendCards()),
        ],
      ),
    );
  }

  Widget _buildSearchTextField() {
    final textField = Padding(
      padding: const EdgeInsets.only(bottom: 20),
      child: TextField(
        controller: _searchTextController,
        autocorrect: false,
        decoration: const InputDecoration(
          hintText: 'Search for a friend',
          prefixIcon: Icon(Icons.search),
        ),
        style: SchejFonts.subtitle,
      ),
    );

    return textField;
  }

  Widget _buildFriendCards() {
    return ListView.builder(
      itemCount: friends.length + 1,
      itemBuilder: (context, index) {
        if (index == friends.length) {
          // Return sized box so FAB doesn't overlap
          return const SizedBox(height: 70);
        }

        final friend = friends[index];
        return Padding(
          padding: const EdgeInsets.only(bottom: 10),
          child: FriendCard(
              name: friend['name'] as String,
              status: friend['status'] as FriendStatus,
              showOverflowMenu: () {},
              curEventName: (friend['curEventName'] ?? '') as String,
              onTap: () => AutoRouter.of(context)
                  .push(FriendSchejPageRoute(name: friend['name'] as String))),
        );
      },
    );
  }
}
