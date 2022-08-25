import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/components/friends/friend_card.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:flutter_app/models/user.dart';
import 'package:flutter_app/router/app_router.gr.dart';
import 'package:provider/provider.dart';

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

    // Start listening to changes.
    _searchTextController.addListener(_updateSearchResults);

    // Initialize the friend search results.
    final api = context.read<ApiService>();
    _updateSearchResults();

    // Add listener to the friends property in ApiService.
    api.addListener(_updateSearchResults, [ApiServiceProperties.friends]);
  }

  void _updateSearchResults() {
    setState(() {
      _friendsSearchResults = context
          .read<ApiService>()
          .getFriendsByQuery(_searchTextController.text);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        _buildSearchTextField(),
        Expanded(child: _buildFriendCards()),
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

  Widget _buildFriendCards() {
    return ListView.builder(
      itemCount: _friendsSearchResults.length + 1,
      itemBuilder: (context, index) {
        if (index == _friendsSearchResults.length) {
          // Return sized box so FAB doesn't overlap
          return const SizedBox(height: 70);
        }

        final friend = _friendsSearchResults[index];
        return Padding(
          padding: SchejConstants.pagePadding
              .copyWith(top: index == 0 ? 10 : 0, bottom: 10),
          child: FriendCard(
              name: friend.fullName,
              picture: friend.picture,
              status: FriendStatus.free /*friend['status'] as FriendStatus*/,
              showOverflowMenu: () {},
              // curEventName: (friend['curEventName'] ?? '') as String,
              onTap: () => AutoRouter.of(context)
                  .push(CompareSchejPageRoute(friendId: friend.id))),
        );
      },
    );
  }
}