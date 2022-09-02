import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/components/friends/friend_card.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:flutter_app/models/status.dart';
import 'package:flutter_app/models/user.dart';
import 'package:flutter_app/router/app_router.gr.dart';
import 'package:provider/provider.dart';

enum _MenuAction {
  compare,
  remove,
}

class FriendsTabWidget extends StatefulWidget {
  const FriendsTabWidget({Key? key}) : super(key: key);

  @override
  State<FriendsTabWidget> createState() => _FriendsTabWidgetState();
}

class _FriendsTabWidgetState extends State<FriendsTabWidget> {
  // Controllers
  final TextEditingController _searchTextController = TextEditingController();
  List<User> _friendsSearchResults = [];

  @override
  void dispose() {
    _searchTextController.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();

    // Initialize the friend search results.
    final api = context.read<ApiService>();
    _updateSearchResults(api);

    // Start listening to changes.
    _searchTextController.addListener(() => _updateSearchResults(api));

    // Add listener to the friends property in ApiService.
    api.addListener(
        () => _updateSearchResults(api), [ApiServiceProperties.friends]);
  }

  void _updateSearchResults(ApiService api) {
    if (!mounted) return;

    // Find and update status for each friend.
    setState(() {
      _friendsSearchResults = api.getFriendsByQuery(_searchTextController.text);
    });
  }

  Future<void> _showMenu(String id, RelativeRect position) async {
    final action = await showMenu<_MenuAction>(
        context: context,
        position: position,
        items: [
          const PopupMenuItem<_MenuAction>(
            value: _MenuAction.compare,
            child: Text('Compare your schej'),
          ),
          const PopupMenuItem<_MenuAction>(
            value: _MenuAction.remove,
            child: Text(
              'Remove friend',
              style: TextStyle(color: SchejColors.red),
            ),
          ),
        ]);

    final api = context.read<ApiService>();
    switch (action) {
      case _MenuAction.compare:
        AutoRouter.of(context).push(CompareSchejPageRoute(
          friendId: id,
          initialIncludeSelf: true,
        ));
        break;
      case _MenuAction.remove:
        api.deleteFriend(id);
        break;
      default:
        break;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        _buildSearchTextField(),
        Expanded(
          child: Consumer<ApiService>(
            builder: (context, api, child) => _buildFriendCards(api),
          ),
        ),
      ],
    );
  }

  Widget _buildSearchTextField() {
    final textField = Padding(
      padding: SchejConstants.pagePadding.copyWith(top: 20, bottom: 10),
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

  Widget _buildFriendCards(ApiService api) {
    return ListView.builder(
      itemCount: _friendsSearchResults.length + 1,
      itemBuilder: (context, index) {
        if (index == _friendsSearchResults.length) {
          // Return sized box so FAB doesn't overlap
          return const SizedBox(height: 70);
        }

        final friend = _friendsSearchResults[index];
        final status = api.friendsStatus[friend.id];

        return Padding(
          padding: SchejConstants.pagePadding
              .copyWith(top: index == 0 ? 10 : 0, bottom: 10),
          child: FriendCard(
              name: friend.fullName,
              picture: friend.picture,
              status:
                  stringToFriendStatus(status != null ? status.status : 'free'),
              curEventName: status?.eventName,
              showOverflowMenu: (RelativeRect position) {
                _showMenu(friend.id, position);
              },
              // curEventName: (friend['curEventName'] ?? '') as String,
              onTap: () => AutoRouter.of(context)
                  .push(CompareSchejPageRoute(friendId: friend.id))),
        );
      },
    );
  }

  FriendStatus stringToFriendStatus(String status) {
    return FriendStatus.values
        .firstWhere((e) => e.toString() == 'FriendStatus.' + status);
  }
}
